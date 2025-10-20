# Alonso Programming Language

**Alonso** is a Formula 1-themed programming language inspired by the legendary Fernando Alonso. This interpreted language transforms traditional programming concepts into exciting F1 racing metaphors, making coding feel like commanding a race car on the circuit.

## Overview

Alonso is a complete programming language featuring:
- **F1-themed syntax** with racing-inspired keywords
- **Tree-walking interpreter** built in Go
- **Dynamic typing** with numbers, strings, booleans, and arrays
- **First-class functions** with lexical scoping
- **Control flow** including loops and conditionals
- **Interactive REPL** for live coding
- **Comprehensive error handling** with detailed messages

## Architecture

The Alonso compiler follows a traditional interpreter architecture:

```
Source Code (.alo) → Lexer → Parser → AST → Interpreter → Output
```

### Core Components

1. **Lexer** (`lexer.go`) - Tokenizes source code into meaningful symbols
2. **Parser** (`parser.go`) - Builds Abstract Syntax Tree (AST) from tokens
3. **AST** (`ast.go`) - Represents program structure as tree nodes
4. **Interpreter** (`interpreter.go`) - Executes AST using tree-walking evaluation
5. **Object System** (`object.go`) - Runtime value representation and environment

## F1-Themed Keywords

| Traditional | Alonso Keyword | Racing Metaphor |
|-------------|----------------|-----------------|
| `var` | `grid` | Starting grid position |
| `func` | `pace` | Racing pace/strategy |
| `if` | `circuit` | Racing circuit decision |
| `else` | `else_circuit` | Alternative racing line |
| `for` | `loop` | Racing loop/lap |
| `while` | `while_racing` | Continue while racing |
| `return` | `return_pit` | Return to pit lane |
| `break` | `break_flag` | Yellow flag (stop) |
| `continue` | `continue_race` | Green flag (continue) |
| `array` | `formation` | Formation lap lineup |

## Quick Start

### Prerequisites
- Go 1.21 or higher

### Installation
```bash
# Clone the repository
git clone <repository-url>
cd alonso

# Build the interpreter
go build -o alonso.exe

# Run interactive REPL
./alonso.exe

# Execute .alo files
./alonso.exe examples/hello.alo
```

## Language Syntax

### Variables (Grid Positions)
```alonso
grid driver_name = "Fernando Alonso"
grid car_number = 14
grid is_racing = true
grid lap_time = 88.5
```

### Functions (Racing Pace)
```alonso
pace calculate_lap_time(base_time, weather_factor) {
    grid adjusted_time = base_time * weather_factor
    return_pit adjusted_time
}

// Function calls
grid result = calculate_lap_time(90.0, 1.1)
```

### Conditionals (Racing Circuits)
```alonso
circuit (weather == "sunny") {
    telemetry("Perfect racing conditions!")
} else_circuit {
    telemetry("Challenging weather")
}

// Complex conditions
circuit (fuel_level > 50 && tire_condition == "good") {
    telemetry("Ready for aggressive strategy")
}
```

### Loops (Racing Laps)
```alonso
// While racing
grid lap = 1
while_racing (lap <= 10) {
    telemetry("Racing lap", lap)
    lap = lap + 1
}

// For loop
loop (grid i = 0; i < 5; i = i + 1) {
    telemetry("Position", i + 1)
}

// Loop control
while_racing (true) {
    circuit (position == 1) {
        break_flag  // Stop racing
    }
    circuit (pit_stop_needed) {
        continue_race  // Skip to next iteration
    }
}
```

### Arrays (Formation)
```alonso
grid drivers = ["Alonso", "Hamilton", "Verstappen"]
telemetry("Pole position:", drivers[0])

// Array operations
grid team_size = length(drivers)
grid updated_team = push(drivers, "Leclerc")
```

## Built-in Functions

- **`telemetry(...)`** - Output function (equivalent to print/console.log)
- **`length(array/string)`** - Returns length of arrays or strings
- **`push(array, element)`** - Adds element to array (returns new array)

## Project Structure

```
alonso/
├── main.go           # Entry point and REPL
├── lexer.go          # Lexical analysis
├── parser.go         # Syntax analysis
├── ast.go            # Abstract Syntax Tree definitions
├── interpreter.go    # Tree-walking interpreter
├── object.go         # Runtime object system
├── examples/         # Sample programs
│   ├── hello.alo
│   ├── functions.alo
│   ├── loops.alo
│   └── conditionals.alo
└── tests/           # Test files
    ├── test.alo
    ├── test_function.alo
    └── minimal.alo
```

## Interactive REPL

The REPL (Read-Eval-Print Loop) provides an interactive environment:

```bash
$ ./alonso.exe
Welcome to Alonso - The F1 Programming Language!
Type 'pit' to exit

alonso> grid speed = 300
alonso> telemetry("Current speed:", speed, "km/h")
Current speed: 300 km/h
alonso> pit
Thanks for racing with Alonso!
```

## Example Programs

### Hello World
```alonso
telemetry("Hello from the Alonso racing circuit!")
grid driver = "Fernando Alonso"
telemetry("Driver:", driver)
```

### Racing Simulation
```alonso
pace simulate_race(laps) {
    grid total_time = 0
    grid current_lap = 1
    
    while_racing (current_lap <= laps) {
        grid lap_time = 88.5 * 1.1
        total_time = total_time + lap_time
        telemetry("Lap", current_lap, "- Time:", lap_time)
        current_lap = current_lap + 1
    }
    
    return_pit total_time
}

grid race_time = simulate_race(5)
telemetry("Total race time:", race_time, "seconds")
```

## Error Handling

Alonso provides comprehensive error messages:
- **Lexical errors** - Invalid characters or unterminated strings
- **Syntax errors** - Malformed expressions or statements
- **Runtime errors** - Type mismatches, undefined variables, division by zero
- **Semantic errors** - Invalid function calls or array access

## Testing

Run the test suite:
```bash
# Test basic functionality
./alonso.exe tests/test.alo

# Test functions
./alonso.exe tests/test_function.alo

# Test minimal program
./alonso.exe tests/minimal.alo
```

## Language Features

### Data Types
- **Numbers** - 64-bit floating point (e.g., `42`, `3.14`)
- **Strings** - UTF-8 text (e.g., `"Fernando Alonso"`)
- **Booleans** - `true` and `false`
- **Arrays** - Dynamic collections (e.g., `[1, 2, 3]`)
- **Functions** - First-class values with closures
- **Null** - Represents absence of value

### Operators
- **Arithmetic** - `+`, `-`, `*`, `/`, `%`
- **Comparison** - `==`, `!=`, `<`, `>`, `<=`, `>=`
- **Logical** - `&&`, `||`, `!`
- **Assignment** - `=`
- **Index** - `array[index]`

### Scoping
- **Lexical scoping** with nested environments
- **Function closures** capture outer scope
- **Block scoping** for control structures

## Future Enhancements

- **Structs/Objects** - Custom data types
- **Modules** - Code organization and imports
- **Standard Library** - Extended built-in functions
- **Optimizations** - Bytecode compilation
- **Debugging** - Step-through debugger
- **Package Manager** - Dependency management

## Contributing

Contributions are welcome! Areas for improvement:
- Additional F1-themed features
- Performance optimizations
- Extended standard library
- Better error messages
- Documentation improvements

## License

This project is open source. Feel free to use, modify, and distribute.

---

*"El Plan is to make programming as exciting as Formula 1 racing!"* 🏎️💨