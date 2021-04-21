load("assert.star", "assert")


def test_count():

    # No args — defaults to 0, 1.
    c0 = count()
    # *step* is omitted when 1. This is equivalent to count(start).
    assert.eq(str(c0), "count(0)")
    assert.eq(next(c0), 0)
    assert.eq(str(c0), "count(1)")
    assert.eq(next(c0), 1)
    assert.eq(str(c0), "count(2)")

    # Only one arg — step defaults to 1.
    c1 = count(5)
    assert.eq(str(c1), "count(5)")
    assert.eq(next(c1), 5)
    assert.eq(str(c1), "count(6)")
    assert.eq(next(c1), 6)
    assert.eq(str(c1), "count(7)")
    assert.eq(next(c1), 7)

    c2 = count(5, 3)
    assert.eq(str(c2), "count(5, 3)")
    assert.eq(next(c2), 5)
    assert.eq(str(c2), "count(8, 3)")
    assert.eq(next(c2), 8)
    assert.eq(str(c2), "count(11, 3)")
    assert.eq(next(c2), 11)

    # Fails
    z = ("a", "b")
    assert.fails(
        lambda: count("a", "b"),
        # fails uses match under the hood, which will use
        # regexp.MatchString, so need to use raw pattern
        # that MatchString would accept.
        r'Got \(\"a\", \"b\"\)',
    )
    assert.fails(
        lambda: count(1, 2, 3),
        r'Got \(1, 2, 3\)'
    )


test_count()
