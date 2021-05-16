load("assert.star", "assert")

def test_main():
    assert.fails(lambda: map(0, []), "want callable")
    assert.fails(lambda: map([], 10), "want callable")
    assert.fails(lambda: map(lambda x: x, 10), "want iterable")
    assert.fails(lambda: map(lambda x: x, [], abc="abc"), "unexpected keyword arguments")

    tests = [
        lambda x: x + 1,
    ]
    for test in tests:
        assert.eq(
            list(map(test, [1, 2, 3, 4])),
            [2, 3, 4, 5],
        )
        assert.eq(
            list(map(test, range(1, 5))),
            [2, 3, 4, 5],
        )

    # Use shortest.
    iter1 = range(1, 5)
    iter2 = range(1, 15)
    iter3 = range(1, 10)
    r = list(map(lambda x, y, z: (x + 1, y + 1, z + 1), iter1, iter2, iter3))
    assert.eq(
        r,
        [(2,2,2), (3,3,3), (4,4,4), (5,5,5),],
    )

    def f(_a):
        return _a

    m = map(f, [1, 2, 3])
    assert.eq(next(m), 1)
    assert.eq(next(m), 2)
    assert.eq(next(m), 3)
    assert.fails(lambda: next(m), "iteration done")
    m = None
    assert.gc()

test_main()