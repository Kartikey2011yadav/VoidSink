# Heffalump Engine

The **Heffalump Engine** is the content generation subsystem of VoidSink. Its primary purpose is to generate an infinite stream of data that mimics legitimate content (like HTML or prose) to keep an attacker's connection open and their parser busy.

## Concept

The engine is based on a **Markov Chain**. A Markov Chain is a stochastic model describing a sequence of possible events in which the probability of each event depends only on the state attained in the previous event.

In the context of text generation:
1.  The engine analyzes a source text (the "corpus").
2.  It calculates the probability of a word appearing after a given prefix (usually 1 or 2 words).
3.  To generate text, it picks a starting word and then randomly selects the next word based on the calculated probabilities.

## Implementation Details

The Heffalump engine in VoidSink is optimized for:
1.  **Speed**: Generation must be faster than the network speed to ensure the buffer is always full.
2.  **Memory Efficiency**: The state map is compact.
3.  **Streaming**: It implements `io.Reader`, allowing it to be piped directly into the network socket.

### The Corpus

VoidSink comes with a default corpus (typically excerpts from public domain literature like *Alice in Wonderland* or *The Metamorphosis*). This ensures the vocabulary is varied and the sentence structure mimics English grammar.

### Why not random bytes?

Sophisticated scanners and bots often have heuristics to detect non-text content. If a web scraper expects HTML and receives random binary noise, it might disconnect immediately. By sending "text-like" garbage, VoidSink tricks the parser into trying to make sense of the data, consuming CPU cycles and keeping the socket open longer.

## Configuration

The engine can be configured to adjust the "order" of the chain (how many previous words influence the next).
- **Order 1**: More random, less coherent.
- **Order 2+**: More coherent, but requires a larger corpus to avoid repetition.

Currently, VoidSink uses a tuned default that balances variety with performance.
