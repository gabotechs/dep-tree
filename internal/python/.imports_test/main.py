import src
import src.main
import src.module
import sys # ignored
from sys import something  # ignored
from src.main import main
import asyncio # ignored
try:
    from .src import main
except Exception as e:
    from .src.main import main
from . import src # should all be imported?
from . import un_existing  # ignored (should error?)
from src.module import *
from src.module import module
from src.module import bar
