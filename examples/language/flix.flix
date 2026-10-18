// Example Flix program written for scc.
// Flix is a functional, imperative, and logic programming language for the JVM.
// See: https://flix.dev

/*
 * Enum representing simple geometric shapes.
 */
enum Shape {
    case Circle(Float64),
    case Square(Float64),
    case Rectangle(Float64, Float64)
}

def area(s: Shape): Float64 = match s {
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
    for (x <- l) yield x * 2

// Imperative example using a mutable reference inside a region.
def countUp(n: Int32): Int32 \ IO =
    region rc {
        let counter = ref 0 @ rc;
        while (deref counter < n and n > 0) {
            counter := deref counter + 1
        };
        deref counter
    }

// A tiny Datalog program embedded directly in Flix.
def reachable(): #{ Path(Int32, Int32) } =
    #{
        Edge(1, 2).
        Edge(2, 3).
        Edge(3, 4).
        Path(x, y) :- Edge(x, y).
        Path(x, z) :- Path(x, y), Edge(y, z).
    }

def main(): Unit \ IO =
    try {
        println("Area of circle: ${area(Shape.Circle(2.0))}");
        println("fib(10) = ${fib(10)}");
        println("Doubled: ${doubled(1 :: 2 :: 3 :: Nil)}");
        println("Counted up to: ${countUp(5)}")
    } catch {
        case ex: Exception => println("Something went wrong: ${ex}")
    }
