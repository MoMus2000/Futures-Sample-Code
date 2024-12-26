package main

import (
  "fmt"
)

type Future interface {
    poll(data interface{}) (interface{}, error)
    then(func (args ...interface{}) Future) Future
}

type FutureDone struct {
    result interface{}
}

func (f *FutureDone) poll(data interface{}) (interface{}, error) {
    return f.result, nil
}

func (f *FutureDone) then(func (args ...interface{}) Future) Future {
    panic("Cannot chain an ended Future")
}

type FutureThen struct{
    left     Future
    right    Future
    callback func (args ...interface{}) Future }

func (f *FutureThen) poll(data interface{}) (interface{}, error){
    if f.left != nil{
      result, _ := f.left.poll(data)
      if result != nil {
        f.right = f.callback(result)
        f.left = nil
        return result, nil
      }
    } else {
        if f.right == nil{
          panic("right side cannot be nil")
        }
        return f.right.poll(data)
    } 
    return nil, nil
}

func (f *FutureThen) then(d (func(args ...interface{}) Future)) Future {
    return &FutureThen{
      left: f,
      right: nil,
      callback: d,
    }
}

type ConcreteFuture struct {
    func_wrapper func(args ...interface{}) interface{}
}

func (cf *ConcreteFuture) poll(data interface{}) (interface{}, error){
    res := cf.func_wrapper(data)
    if res == nil{
      return false, nil
    }
    return true, nil
}

func (cf *ConcreteFuture) then(data func (args ...interface{}) (Future)) Future {
    return &FutureThen{
        left: cf,
        right: nil,
        callback: data,
    }
}

type Scheduler struct {
    futures []Future
}

func (s *Scheduler) add_future(f Future) {
    s.futures = append(s.futures, f)
}

func (s *Scheduler) start(){
  // for {
    for _, future := range s.futures {
        future.poll(nil)
    }
  // }
}

func counter(args ...interface{}) interface{}{
  for i:=0; i< 10; i++{
    fmt.Printf("Counter @ %d\n", i);
  }
  return nil
}

func main(){
  var cf = ConcreteFuture{
    func_wrapper: counter,
  }
  var next = cf.then(
    func(args ...interface{}) Future {
      fmt.Println("Calling the callback for next")
      return &FutureDone{result:"done"}
    },
  )
  var one_after = next.then(
    func(args ...interface{}) Future {
      fmt.Println("Calling the callback for one_after")
      return &FutureDone{result:"done"}
    },
  )
  var sf = ConcreteFuture{
    func_wrapper: counter,
  }
  var nxt = sf.then(
    func(args ...interface{}) Future {
      fmt.Println("Calling the callback for nxt")
      return &FutureDone{result:"done"}
    },
  )
  var ne_after = nxt.then(
    func(args ...interface{}) Future {
      fmt.Println("Calling the callback for ne_after")
      return &FutureDone{result:"done"}
    },
  )

  scheduler := Scheduler{}
  scheduler.add_future(one_after)
  scheduler.add_future(ne_after)
  scheduler.start()

}
