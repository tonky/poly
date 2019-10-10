#[cfg(test)]
mod tests {
    use reqwest;
    use scopeguard::defer;
    use std::process::{Command,Child};
    use std::{thread, time};


    fn init_and_run_on_cluster() -> Child {
        let sec = time::Duration::from_millis(1000);

        let child = Command::new("skaffold")
                    .arg("dev")
                    .arg("--port-forward")
                    .current_dir("..")
                    .spawn()
                    .expect("failed to execute process");

        loop {
            let resp = reqwest::get("http://localhost:9000/");
            // println!("{:#?}", resp);

            match resp {
                Ok(r) => { println!("ok resp: {:?}", r); break }
                Err(e) => { println!("not ready yet, sleeping: {}", e); thread::sleep(sec); continue }
            }
        }

        println!("Init tests: done!");
        child
    }

    #[test]
    fn get_store() -> Result<(), Box<dyn std::error::Error>> {
        let mut child = init_and_run_on_cluster();

        {
            defer! {{
                println!("Killing child");
                child.kill().expect("command wasn't running");
            }};


            let resp = reqwest::get("http://localhost:9000/store/")?.text()?;
            println!("{:#?}", resp);

            assert_eq!(resp, "Hello store!\n");
        }

        Ok(())
    }
}
