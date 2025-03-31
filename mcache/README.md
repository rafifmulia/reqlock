# Description
This package (mcache) implement hash table cache.
This package is not suit for anyone who wants to store value in memory (being cached),
because this package will turns your cache into hash value.
```map[value any]bool```
That means you cannot retrieve the real value of cache, just retrieves the status of cache by its value.

# How to use
If you have concurrent request called "bookseat",
you can store all values into cache with `mcache.Set("bookseat", value any)`.
Then you must write code like these before server runs:
```
go mcache.CleanupRoutine(1*time.Minute, 50) // Runs every 1 minute and clean expired cache more than 50 seconds.
defer ShutdownCleanupRoutine()
```