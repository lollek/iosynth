#! /bin/bash
for o in {2..5}; do
    for n in C c D d E F f G g A a B; do
        echo "$o$n"
        echo -en "2$o$n" | nc -w0 -u4 localhost 49161
        sleep "0.1"
    done
done
