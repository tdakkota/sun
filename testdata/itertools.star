load("assert.star", "assert")


def test_count():

    c0 = count(0, 1)
    # *step* is omitted when 1. This is equivalent to count(start).
    assert.eq(str(c0), "count(0)")
    assert.eq(next(c0), 0)
    assert.eq(str(c0), "count(1)")
    assert.eq(next(c0), 1)
    assert.eq(str(c0), "count(2)")

    c1 = count(0, 5)
    assert.eq(str(c1), "count(0, 5)")
    assert.eq(next(c1), 0)
    assert.eq(str(c1), "count(5, 5)")
    assert.eq(next(c1), 5)
    assert.eq(str(c1), "count(10, 5)")
    assert.eq(next(c1), 10)

    c2 = count(5, 3)
    assert.eq(str(c2), "count(5, 3)")
    assert.eq(next(c2), 5)
    assert.eq(str(c2), "count(8, 3)")
    assert.eq(next(c2), 8)
    assert.eq(str(c2), "count(11, 3)")
    assert.eq(next(c2), 11)


test_count()
