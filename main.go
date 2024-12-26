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

func NewConcreteFuture(func_wrapper func(args ...interface{}) interface{}) *ConcreteFuture{
  return &ConcreteFuture{
    func_wrapper,
  }
}

func (cf *ConcreteFuture) poll(data interface{}) (interface{}, error){
    res := cf.func_wrapper(data)
    return res, nil
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
    for _, future := range s.futures {
        go future.poll(nil)
    }
    fmt.Println("Executed ...")
    for {

    }
}

func counter(args ...interface{}) interface{}{
  for i:=0; i< 10; i++{
    fmt.Printf("Counter @ %d\n", i);
  }
  return true
}

type CounterFuture struct {
    start int
    end   int
}

func NewCounterFuture(start int, end int) *CounterFuture {
  return &CounterFuture{
      start,
      end,
  }
}

func (c *CounterFuture) poll(data interface{}) (interface{}, error) {
    if c.start < c.end {
      c.start += 1
      fmt.Printf("Polling %d \n", c.start)
      return nil, nil
    }
    return c.start, nil
}

func (c *CounterFuture) then(data func (args ...interface{}) (Future)) Future {
    return &FutureThen{
        left: c,
        right: nil,
        callback: data,
    }
}

func main(){
  var cf = NewConcreteFuture(counter).then(
    func(args ...interface{}) Future {
      fmt.Println("cf 1")
      return &FutureDone{result:"done"}
    },
  ).then(
    func(args ...interface{}) Future {
      fmt.Println("cf 2")
      return &FutureDone{result:"done"}
    },

  )

  var sf = NewConcreteFuture(counter).then(
    func(args ...interface{}) Future {
      fmt.Println("sf 1")
      return &FutureDone{result:"done"}
    },
  ).then(
    func(args ...interface{}) Future {
      fmt.Println("sf 2")
      return &FutureDone{result:"done"}
    },
  )

  scheduler := Scheduler{}
  scheduler.add_future(cf)
  scheduler.add_future(sf)
  scheduler.start()

}
