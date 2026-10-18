// Example Flix program written for scc.
// Flix is a functional, imperative, and logic programming language for the JVM.
// See: https://flix.dev

import java.lang.Exception

/*
 * Enum representing simple geometric shapes.
 * /* Block comments in Flix nest, unlike C or Java. */
 * This whole comment still counts because nesting keeps it open.
 */
enum Shape {
    case Circle(Float64),
    case Square(Float64),
    case Rectangle(Float64, Float64)
}

def area(shape: Shape): Float64 = match shape {
    case Shape.Circle(r)       => 3.14159 * r * r
    case Shape.Square(s)       => s * s
    case Shape.Rectangle(w, h) => w * h
}

def fib(n: Int32): Int32 =
    if (n < 2)
        n
    else
        fib(n - 1) + fib(n - 2)

def sumEvens(l: List[Int32]): Int32 =
    List.foldLeft((acc, x) -> if (Int32.remainder(x, 2) == 0) acc + x else acc, 0, l)

def doubled(l: List[Int32]): List[Int32] =
    List.map(x -> x * 2, l)

// Imperative example using a mutable reference inside a region.
def countUp(n: Int32): Int32 =
    region rc {
        let counter = Ref.fresh(rc, 0);
        foreach (_ <- List.range(0, n))
            Ref.put(Ref.get(counter) + 1, counter);
        Ref.get(counter)
    }

// A tiny Datalog program embedded directly in Flix, computing transitive
// reachability over a fixed set of edges.
def reachable(): Vector[(Int32, Int32)] =
    let pr = #{
        Edge(1, 2).
        Edge(2, 3).
        Edge(3, 4).
        Path(x, y) :- Edge(x, y).
        Path(x, z) :- Path(x, y), Edge(y, z).
    };
    query pr select (x, y) from Path(x, y)

def main(): Unit \ IO =
    try {
        println("Area of circle: ${area(Shape.Circle(2.0))}");
        println("fib(10) = ${fib(10)}");
        println("Doubled: ${doubled(1 :: 2 :: 3 :: Nil)}");
        println("Counted up to: ${countUp(5)}");
        println("Reachable pairs: ${reachable()}")
    } catch {
        case _: Exception => println("Something went wrong")
    }
