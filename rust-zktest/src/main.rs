use r1cs::{Bn128, GadgetBuilder, Expression, Element, values};

fn main() {
    println!("Hello, world!");

    // Create a gadget which takes a single input, x, and computes x*x*x.
    let mut builder = GadgetBuilder::<Bn128>::new();
    let x = builder.wire();
    println!("r1cs: {}", x);

    let x_exp = Expression::from(x);
    println!("x_exp: {}", x_exp);
    let x_squared = builder.product(&x_exp, &x_exp);
    println!("x_squared: {}", x_squared);
    let x_cubed = builder.product(&x_squared, &x_exp);
    println!("x_cubed: {}", x_cubed);
    let gadget = builder.build();
    // println!("gadget: {}", gadget);

    let mut values = values!(x => 5u8.into());
    // println!("values: {}", values);

    let constraints_satisfied = gadget.execute(&mut values);
    assert!(constraints_satisfied);
    println!("constraints_satisfied: {}", constraints_satisfied);

    assert_eq!(Element::from(125u8), x_cubed.evaluate(&values));

}
