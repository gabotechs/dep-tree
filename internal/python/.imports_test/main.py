import src
import src.main
import src.module
import sys
from sys import something
from src.main import main
import asyncio
try:
    from .src import main
except Exception as e:
    from .src.main import main
from . import src
from . import un_existing
from src.module import *
from src.module import module
from src.module import bar
