use fern;
use tracing::{warn, info, trace, debug, error, span, Level, instrument};
use chrono;
use tracing_subscriber::FmtSubscriber;

/*
fn setup_logger() -> Result<(), fern::InitError> {
    fern::Dispatch::new()
        .format(|out, message, record| {
            out.finish(format_args!(
                "{}[{}][{}] {}",
                chrono::Local::now().format("[%Y-%m-%d][%H:%M:%S]"),
                record.target(),
                record.level(),
                message
            ))
        })
        .level(Level::DEBUG)
        .chain(std::io::stdout())
        .chain(fern::log_file("output.log")?)
        .apply()?;
    Ok(())
}
*/
#[cfg(test)]
mod tests {
    use reqwest;
    use scopeguard::defer;
    use std::process::{Command,Child,Stdio};
    use std::{thread, time};
    use std::io::Read;
    use std::rc::Rc;
    use super::*;

    #[instrument]
    fn init_and_run_on_cluster() -> Child {
        // setup_logger().unwrap();
        let my_subscriber = FmtSubscriber::new();

        tracing::subscriber::set_global_default(my_subscriber).expect("setting tracing default failed");

        let sec = time::Duration::from_millis(500);

        loop {
            match Command::new("skaffold").arg("run").current_dir("..").status() {
                Ok(status) => {
                    info!("skaffold ready: {:?}", status);
                    if status.success() { break };
                }
                Err(e) => { info!("skaffold not ready yet, sleeping: {}", e); }
            }

            thread::sleep(sec);
        };

        let child = loop {
            // kubectl port-forward service/gw-lb 8888:8080
            match Command::new("kubectl")
                .arg("port-forward")
                .arg("service/api-gateway")
                .arg("9000:8080")
                .stdout(Stdio::piped())
                // .stderr(Stdio::piped())
                .spawn() {
                    Ok(mut child) => {
                        info!("ok kubectl child: {:?}", child);

                        let mut buffer = [0; 10];

                        let _bytes_read = child.stdout.as_mut().expect("no stdout?").read_exact(&mut buffer);

                        let got = String::from_utf8_lossy(&buffer);

                        info!("got buffer: {}: '{}'", got.len(), got);

                        if got.contains("Forwarding") {
                            info!("okay, no errors on proxying!");

                            break child;
                        } else {
                            warn!("error proxying, buffer: '{}'", got);
                            child.kill();
                        }
                    }

                    Err(e) => { error!("kubectl port-forward spawn error: {}", e); }
                }

            thread::sleep(sec);
        };

        loop {
            let resp = reqwest::get("http://localhost:9000/");
            // println!("{:#?}", resp);

            match resp {
                Ok(r) => { println!("ok resp: {:?}", r); break }
                Err(e) => { println!("not ready yet, sleeping: {}", e) }
            }

            thread::sleep(sec);
        }

        println!("Init tests: done!");
        child
    }

    #[test]
    #[instrument]
    fn get_store() -> Result<(), Box<dyn std::error::Error>> {
        let mut child = init_and_run_on_cluster();

        {
            defer! {{
                info!("Killing child...");

                child.kill().expect("command wasn't running");

                let status = Command::new("skaffold")
                            .arg("delete")
                            .current_dir("..")
                            .status()
                            .expect("failed to execute process");

                assert!(status.success());
            }};


            let client = reqwest::Client::builder()
                .timeout(time::Duration::from_secs(1))
                .build()?;

            let resp = client.get("http://localhost:9000/store/").send()?.text()?;

            info!("Gor stores resp: {:#?}", resp);

            assert_eq!(resp, "Hello store!\n");
        }

        Ok(())
    }
}
