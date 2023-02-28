mod sum;
mod div;
mod avg;
mod abs;
mod avg_2;

pub use crate::abs::abs::abs;
pub use crate::div::div;
pub use crate::avg_2::avg::avg;
pub use self::sum;
pub use self::sum::*;

pub fn run() {
    sum::sum();
    div();
    abs();
    avg();
    avg::avg();
}

#[cfg(test)]
mod tests {
    use crate::run;

    #[test]
    fn it_works() {
        run()
    }
}
