xml-dom
=======

DOM implementation for Golang.

The goal is to have a DOM implementation that conforms to the W3C DOM Recommendation. For the moment only part of the DOM Level 1 is implemented, with almost no namespace support.

A side goal is to have a DOM implementation that uses the Golang XML parser and that is able to output the same document with as little change as necessary. As such, it keeps insignificant whitespace inside DOM elements such that the output can be byte to byte equal to the input, unless you change the DOM.
