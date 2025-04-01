# Description
This package (mcache) implement hash table cache.
This package is not suit for anyone who wants to store value in memory (being cached),
because this package will turns your cache into hash value `map[value any]int64`.
That means you cannot retrieve the real value of cache, just retrieves the status of cache by its value.
I can say this package is concurrent-safe, because the test case when i run was 100% PASS,
but i don't guarantee 1000% its concurrent-safe, maybe i made a mistake.

# Data Structure
cacheKey -> The concurrent request name.
cacheVal -> Cache value (being hashed) and last modified time in Unix time seconds.

# How to use
If you have concurrent request called "bookseat",
you can store all values into cache with `mcache.Set("bookseat", value any)`.
Then you must write code like these before the server runs:
```
go mcache.CleanupRoutine(1*time.Minute, 50) // Runs every 1 minute and clear staled cache more than 50 seconds.
defer mcache.ShutdownCleanupRoutine()
```