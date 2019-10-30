[![Build Status](https://travis-ci.com/quapka/PA193_mnemonic_aller.svg?branch=master)](https://travis-ci.com/quapka/PA193_mnemonic_aller)
[![Coverage Status](https://coveralls.io/repos/github/quapka/PA193_mnemonic_aller/badge.svg?branch=master)](https://coveralls.io/github/quapka/PA193_mnemonic_aller?branch=master)

PA193_mnemonic_aller
====================

Language: Go


TODO
====

Testing:
- check memory requirements (run with restricted memory)
- reason about side-channels and constant-time algorithms (is there any leakage?)
- (TDD) unit test everything (extreme values - min, max, pattern values)
- check inputs provided by the users (can they produce bug/error - e.g. end up in an error message?)
- go through this [list](https://github.com/mre/awesome-static-analysis#go) and maybe include in out CI
- (check for change in the size of the binary)
