from abc import ABC, abstractmethod
from typing import Callable
import copy

class Future(ABC):
    @abstractmethod
    def poll(self, data):
        pass

class FutureDone(Future):
    def __init__(self, result):
        self.result = result

    def poll(self, data):
        return self.result

class FutureThen(Future):
    def __init__(self, left, right, func):
        self.left = left
        self.right = right
        self.func = func

    def poll(self, data):
        if self.left is not None:
            result = self.left.poll(data)
            if result is not None:
                print(f"FutureThen resolved with result: {result}")
                self.right = self.func(result)
                self.left = None
            return None
        else:
            assert self.right is not None, "Right side should not be none"
            return self.right.poll(data)

    def then(self, func: Callable):
        print("Going for the second chain")
        return copy.deepcopy(FutureThen(left=self, right=None, func=func))

def some_sample_then_func(result):
    print("In `some_sample_then_func`: Received", result)
    return FutureDone("Donzo")

def some_sample_then_func2(result):
    print("In `some_sample_then_func2`: Received", result)
    return FutureDone("Donzo")

class CounterFuture(Future):
    def __init__(self, end, start=0):
        self.start = start
        self.end = end

    def poll(self, data):
        if self.start < self.end:
            print(f"CounterFuture polling: {self.start}")
            self.start += 1
            return None
        print(f"CounterFuture completed at {self.start}")
        return self.start

    def then(self, func: Callable):
        print("Doing the first chain")
        print(self)
        print(func)
        return copy.deepcopy(FutureThen(left=self, right=None, func=func))

if __name__ == "__main__":
    # Chain futures

    count = CounterFuture(10).then(
        lambda _:  FutureDone("Done BOI")
    ).then(
        lambda _ : FutureDone("done for good")
    )

    count2 = CounterFuture(start=10, end=50).then(
        lambda _:  FutureDone("Done BOI")
    ).then(
        lambda _ : FutureDone("done for good")
    )

    futures = [
        count,
        count2
    ]

    print("Chained Future Created")

    # Poll the future
    while True:
        for future in futures:
            future.poll(None)

