# generationk

The inspiration for this project took place after using a few other backtesting frameworks in Python. I was tired of waiting for results and concluded that I want the fast feeling of a compiled language and I also want to make use the multiple processor cores that often is available but rarely used.

## Design choices

This is is really tricky for me and I have started over a few times. I backtest mostly on daily data so I have questioned is there even a need for a backtesting framwork, I could just work with an alternative for pandas. Is there then really a point for event driven, I can just approximate slippage and comissions, its not a big deal for my small trading. Going for event driven in the end came more from a point of being able to split the program up to run on multiple computers for performance reasons rather than a needs for realism (but implictly contributes to realism).

Another choice coming from an object oriented background is the inheritance and encapsulation of data where things magically happens vs. seing a lot of float arrays and increasing counters yourself in the strategy. From the start I wanted as a simple approach as possible, i.e. working with float arrays inside the strategy but coming from the object oriented world I tend to complicate things that way.

## The Crossing MA example looks like this

