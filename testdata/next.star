load("assert.star", "assert")

def test_main():
    assert.fails(lambda: next([]), "iteration done")
    assert.eq(
        next([1, 2, 3]),
        1,
    )
    assert.eq(
        next([], 1),
        1,
    )

test_main()