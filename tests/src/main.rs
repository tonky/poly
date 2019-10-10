fn main() {
    println!("Hello, world 2!");
}

#[cfg(test)]
mod tests {
    #[test]
    fn test_add() {
        assert_eq!(1+2, 4);
    }
}
