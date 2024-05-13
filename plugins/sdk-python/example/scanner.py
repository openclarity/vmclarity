#!/usr/bin/env python3

from plugin.scanner import AbstractScanner
from plugin.server import Server


# implement abstract class
class ExampleScanner(AbstractScanner):
    def __init__(self):
        return


if __name__ == '__main__':
    Server.run(ExampleScanner())
