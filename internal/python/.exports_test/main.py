import foo as foo_2
import folder.foo as foo_3
from foo import bar
from folder.foo import baz as baz_2
from .src import a
from .module import *
from .module import module


foo = 1

foo_1, foo_2 = 1, 2

(
    foo_3,
    foo_4
) = 1, 2

foo_5 = foo_6 = 5


def func(a, b, c):
    pass

class Class(Other):
    pass
