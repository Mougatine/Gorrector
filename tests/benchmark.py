#! /usr/bin/env python3

import argparse
import sys
import time
import os
import subprocess

def format_query(words, distance):
    prefix = 'approx {} '.format(str(distance))
    queries = ['{}{}'.format(prefix, word) for word in words]
    return '\n'.join(queries)


def read_n_words(path, nb):
    words = []
    parse = lambda x: x.split()[0]
    with open(path, 'r') as f:
        for _ in range(nb):
            words.append(parse(f.readline()))
            """
            for _ in range(0, 1000):
                f.readline()
            """

    return words


def run(args, words):
    cmd = 'echo "{}" | ./{} {}'.format(words, args.app, args.trie)
    ps = subprocess.Popen(cmd, shell=True,
                          stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    ps.wait()


def logic(argv):
    parser = argparse.ArgumentParser(description='Benchmarker.')
    parser.add_argument('--app', action='store', dest='app', required=True)
    parser.add_argument('--trie', action='store', dest='trie', required=True)
    parser.add_argument('--run', action='store', type=int, default=10, dest='run')
    parser.add_argument('--words', action='store', dest='words', required=True)
    parser.add_argument('--dist', action='store', dest='dist', type=int, default=0)

    args = parser.parse_args(argv)
    words = read_n_words(args.words, args.run)
    words = format_query(words, args.dist)

    begin = time.time()
    run(args, words)
    diff = float(time.time() - begin)
    print('Run {} queries in {}s.'.format(args.run, diff))
    print('Thus {} query par second.'.format(args.run / diff))


if __name__ == '__main__':
    logic(sys.argv[1:])
    