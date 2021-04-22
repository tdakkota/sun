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

    # Negative args.
    c3 = count(-5, -10)
    assert.eq(str(c3), "count(-5, -10)")
    assert.eq(next(c3), -5)
    assert.eq(str(c3), "count(-15, -10)")
    assert.eq(next(c3), -15)

    c4 = count(5, -5)
    assert.eq(str(c4), "count(5, -5)")
    assert.eq(next(c4), 5)
    assert.eq(str(c4), "count(0, -5)")
    assert.eq(next(c4), 0)
    assert.eq(str(c4), "count(-5, -5)")
    assert.eq(next(c4), -5)

    # Int start, float step.
    c5 = count(0, 0.1)
    assert.eq(str(c5), "count(0, 0.1)")
    assert.eq(next(c5), 0)
    assert.eq(str(c5), "count(0.1, 0.1)")
    assert.eq(next(c5), 0.1)
    assert.eq(str(c5), "count(0.2, 0.1)")
    assert.eq(next(c5), 0.2)

    # Float start, int step — this should be handled same as above
    # but check to be exhaustive.
    c6 = count(0.5, 5)
    assert.eq(str(c6), "count(0.5, 5)")
    assert.eq(next(c6), 0.5)
    assert.eq(str(c6), "count(5.5, 5)")
    assert.eq(next(c6), 5.5)
    assert.eq(str(c6), "count(10.5, 5)")
    assert.eq(next(c6), 10.5)

    # This test may seem similar to c5 but is different because
    # here step > 1. In the case that 0 < step < 1, fmt.Sprintf,
    # which is used in String(), will display it as a float but
    # may display it as an int if the proper flags aren't used.
    c7 = count(5.0, 0.5)
    assert.eq(str(c7), "count(5.0, 0.5)")
    assert.eq(next(c7), 5.0)
    assert.eq(str(c7), "count(5.5, 0.5)")
    assert.eq(next(c7), 5.5)
    assert.eq(str(c7), "count(6.0, 0.5)")
    assert.eq(next(c7), 6.0)

    # NaNs
    c8 = count(0, float('nan'))
    assert.eq(str(c8), "count(0, %s)" % (float('nan')))
    assert.eq(next(c8), 0)
    assert.eq(str(c8), "count(%s, %s)" % (float('nan'), float('nan')))
    assert.eq(next(c8), float('nan'))
    assert.eq(str(c8), "count(%s, %s)" % (float('nan'), float('nan')))
    assert.eq(next(c8), float('nan'))

    c9 = count(0, float("+inf"))
    assert.eq(str(c9), "count(0, %s)" % (float("+inf")))
    assert.eq(next(c9), 0)
    assert.eq(str(c9), "count(%s, %s)" % (float("+inf"), float("+inf")))
    assert.eq(next(c9), float("+inf"))
    assert.eq(str(c9), "count(%s, %s)" % (float("+inf"), float("+inf")))
    assert.eq(next(c9), float("+inf"))

    c10 = count(0, float("-inf"))
    assert.eq(str(c10), "count(0, %s)" % (float("-inf")))
    assert.eq(next(c10), 0)
    assert.eq(str(c10), "count(%s, %s)" % (float("-inf"), float("-inf")))
    assert.eq(next(c10), float("-inf"))
    assert.eq(str(c10), "count(%s, %s)" % (float("-inf"), float("-inf")))
    assert.eq(next(c10), float("-inf"))

    c11 = count(float("nan"), 2)
    assert.eq(str(c11), "count(%s, 2)" % (float('nan')))
    assert.eq(next(c11), float('nan'))
    assert.eq(str(c11), "count(%s, 2)" % (float('nan')))
    assert.eq(next(c11), float('nan'))

    c12 = count(float("+inf"), 2)
    assert.eq(str(c12), "count(%s, 2)" % (float('+inf')))
    assert.eq(next(c12), float('+inf'))
    assert.eq(str(c12), "count(%s, 2)" % (float('+inf')))
    assert.eq(next(c12), float('+inf'))

    c13 = count(float("-inf"), 2)
    assert.eq(str(c13), "count(%s, 2)" % (float('-inf')))
    assert.eq(next(c13), float('-inf'))
    assert.eq(str(c13), "count(%s, 2)" % (float('-inf')))
    assert.eq(next(c13), float('-inf'))

    # Fails
    # Non-numeric arg fails.
    assert.fails(
        lambda: count("a", "b"),
        # fails uses match under the hood, which will use
        # regexp.MatchString, so need to use raw pattern
        # that MatchString would accept.
        r'Got \(\"a\", \"b\"\)',
    )

    # Too many arg fails — should be handled by UnpackArgs but
    # check to be exhaustive.
    assert.fails(
        lambda: count(1, 2, 3),
        r'Got \(1, 2, 3\)'
    )


test_count()
