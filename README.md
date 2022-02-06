# miniscala
miniscala programming language

Voici un exemple of the program written in miniscala

```Scala
  
  def fac(x: Int): Int{
      if (x == 1) {
          return 1
      }
      return x * fac(x - 1)
  }


  def fib(n: Int): Int {
      if (n <= 1) {
          return 1
      }
      return fib(n - 1) + fib(n - 2)
  }

  print(fac(fib(5))) // outputs 40320
```

For more examples, refer to source files located under the "sources" folder. 
