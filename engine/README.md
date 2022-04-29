

# Search Benchmarking

Synopsis:

* search functions do a single search for a fixed Board object (`boardForBenchMark`)
* we use Go's benchmarking tool in combination with our method `engine.SingleSearch()`
* fixed depth that is feasable for all search functions

| method / changes                        | depth |       # nodes |   b.N |           ns/op |     allocs/op |          B/op |
| --------------------------------------- | ------| ------------- | ----- | --------------- | ------------- | ------------- |
| `AB1` (`searchAlphaBetaGreedy`)^1       |    10 |        32,242 |   181 |       6,096,807 |        44,938 |     2,073,042 |
| `AB2` (`searchAlphaBeta`, non-greedy)   |    10 |        17,801 |   477 |       2,589,180 |        20,234 |       323,748 |
| `board.New()`: add capacity             |    10 |        17,801 |   481 |       2,577,734 |        20,234 |       323,744 |
| `board.DoMovePro()`: omit score         |    10 |        17,801 |   510 |       2,318,421 |        20,234 |       323,744 |
| experimental `board.DoLeftProRow()`^2   |    10 |       331,327 |   31  |      37,679,491 |       364,729 |     5,835,664 |
| -> "normalized" to match previous bench |    -- |      "17,801" |    -- |       2,024,383 |        19,596 |       313,529 |
Foornotes:

1. Number of Nodes is higher than for `AB2` because more children are created than we need to search into.
2. Experimented to have no loops in the `b.LeftProPro()` function. For each (fixed) row we just call
`LeftProRow()` with pointers to the slice elements. The latter is inaccurate and tests are failing.

| method / changes                  | depth |       # nodes |   b.N |           ns/op |     allocs/op |          B/op |
| ----------------------------------| ------| ------------- | ----- | --------------- | ------------- | ------------- |
| `AB3`                             |    15 |     1,002,906 |    10 |     114,665,080 |     1,063,801 |    17,021,030 |
| `AB3`                             |    25 | 1,253,504,795 |     1 | 137,146,747,021 | 1,342,669,159 | 1,342,669,159 |

## SearchTarget
some simple benchmarking of SearchTarget():

| search function        | evaluation     | target     | time      | nodes      | maxDepth | depth | notes                             |
| ---------------------- | -------------- | ---------- | --------- | ---------- | -------- | ----- | --------------------------------- |
| searchTarget -> \*Node | hasTarget      | 32         | 23.9s     | 6,151,755  | 40       | 37    | naive search for target node      |

Note: *There was a bug in the ExpandSeaarch function, "enabling" to search much deeper than we would.*


| search function        | evaluation     | target     | time      | nodes      | maxDepth | depth | notes                             |
| ---------------------- | -------------- | ---------- | --------- | ---------- | -------- | ----- | --------------------------------- |
| searchTarget -> \*Node | hasTarget      | 8          | 16.37s    | 51,233,399 | 10       | 3     | fixed                             |

Note: *not able to find any target higher than 8*.

## Alpha Beta pruning
Improved performance a lot, but not benchmarked. This allows us to search the tree partially, 
i.e. we have time to search deeper.

At the same time, another improvement was not to use *board.GameOverX()* as a termination condition. 
This function is quite expensive and does work that is already done by *_spexpand()*. So instead we 
check the number of children returned by the latter (if *len(children) == 0* => game over).

## Adaptive depth
Depth is adapted each time it's player's turn. See *AdjustDepth()*. Note that from now on instead
of "MaxDepth", we use a parameter/flag for default depth.

## Greedy vs. lazy search expansion
Previously used *_spexpand()*, which creats and returns a list of boards. This is unefficient since:
* many of these boards are not needed (alpha-beta pruning)
* they need to stored (=copied to) memory until the recursive call to the current child returns

So instead we use *nextChildForPlayer()* and *nextChildForOpponent()*. We only generate one child 
at a time, and only variable we need to store is *nextIndex* (and maybe *childCount*).

Comparison of the two methods, with **defaultDepth=7** and **speval()** for evaluation. Average 
think time (seconds):

- greedy generation: 0.1636
- lazy generation:   0.0782

Not a very solid "benchmark" (only run once), but suggests more than 200% speed-up.

## Evaluation improvement
- Tried to replace *speval()* with *speval2()* which does not division and returns uint8. 
But game performance dropped down dramatically.

- Tried to replace *speval()* with *speval3()* which returns only the count of zeros, uint8 (i.e. 
rewards keeping many tiles empty). Game performance did not seem to change. Time performance 
dropped (!!!): **average think time (s): 0.08146** (defaultDepth=7).


## 25.11.2021 - Goroutines
Use goroutines for the top-level search call. I.e. use 4 goroutines, instead of 1. Does not seem to
improve time performance.

## 25.11.2021 - other improvement attempts:
- board.New(): add capacity equal to params.N (maybe slight, ~5%, improvement in time, but MAYBE).


## 26.11.2021 - benchmarking
* Cleaned up code a little bit, moved dead code to `engine/arxv/`
* A more reliable benchmarking format. See section Benchmark