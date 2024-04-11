from plugin.scanner.server import Server
from scanner import ExampleScanner

if __name__ == '__main__':
    scanner = ExampleScanner()
    server = Server(scanner)
    server.start(host="0.0.0.0", port=8080)
