import foo as foo_2
import folder.foo as foo_3

try:
  import foo
except:
  import folder.foo as foo

try:
  from foo import bar # comment
except:
  from folder.foo import baz as baz_2

# from .src import a
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
    """A prompt template for chat models.

    Use to create flexible templated prompts for chat models.

    Examples:
        .. code-block:: python
            from . import main
            from langchain.prompts import ChatPromptTemplate

            template = ChatPromptTemplate.from_messages([
                ("system", "You are a helpful AI bot. Your name is {name}."),
                ("human", "Hello, how are you doing?"),
                ("ai", "I'm doing well, thanks!"),
                ("human", "{user_input}"),
            ])

            messages = template.format_messages(
                name="Bob",
                user_input="What is your name?"
            )
    """
    pass

class Class(Other):
    """Get a new ChatPromptTemplate with some input variables already filled in.

    Args:
        **kwargs: keyword arguments to use for filling in template variables. Ought
                    to be a subset of the input variables.

    Returns:
        A new ChatPromptTemplate.


    Example:

        .. code-block:: python
            from . import main
            from langchain.prompts import ChatPromptTemplate

            template = ChatPromptTemplate.from_messages(
                [
                    ("system", "You are an AI assistant named {name}."),
                    ("human", "Hi I'm {user}"),
                    ("ai", "Hi there, {user}, I'm {name}."),
                    ("human", "{input}"),
                ]
            )
            template2 = template.partial(user="Lucy", name="R2D2")

            template2.format_messages(input="hello")
    """
    pass

try:
  # This import only works on python 3.3 and above.
  import collections.abc as collections_abc  # pylint: disable=unused-import
except ImportError:
  import collections as collections_abc  # pylint: disable=unused-import

from foo import (
    a,
    b,
    c
)
