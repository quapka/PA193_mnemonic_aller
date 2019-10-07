#!/usr/bin/python3.6

import hashlib
import os

def bit_len(x):
    if type(x) != int:
        raise ValueError("'x' needs to be an integer")
    return len(bin(abs(x))[2:])

def mnem(x_bytes, wordlist):
    # check if input is bytes
    if type(x_bytes) != bytes:
        raise ValueError("Entropy is not 'bytes'")
    # check if input is multiple of 32
    x_hexa = x_bytes.hex()
    ENT = bit_len(int(x_hexa, 16))
    if ENT % 32 != 0:
        raise ValueError('Entropy is not a multiple of 32')
    # check if input is from the proper range
    if ENT < 128 or ENT > 256:
        raise ValueError('Length entropy is not from the accepted range')
    # calculate the check sum
    m = hashlib.sha256()
    m.update(x_bytes)
    x_digest = m.digest()
    checksum_len = ENT / 32
    checksum_bits = bin(int(x_digest.hex(), 16))[2: int(2 + checksum_len)]
    x_bytes += bytes([int(checksum_bits, 2)])
    # create groups by eleven
    x_bits = bin(int(x_bytes.hex(), 16))[2:]
    groups = [x_bits[11 * i: 11 * (i + 1)] for i in range(len(x_bits) // 11)]
    for g in groups:
        print(wordlist[int(g, 2)])


if __name__ == '__main__':
    with open('../wordlists/english.txt', 'r') as f:
        WORDLIST = [x.strip() for x in f.readlines()]

    # generate exactly 128 bits of entropy
    entropy = os.urandom(16)
    while bit_len(int(entropy.hex(), 16)) != 128:
        entropy = os.urandom(16)

    mnem(entropy, WORDLIST)
