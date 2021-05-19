load("assert.star", "assert")


def test_count():

    ### Hash ###
    assert.true(count(0) != count(0))
    assert.true(count(5, 0.6) != count(5, 0.6))
    h_c0 = count(1, 2)
    assert.true(h_c0 == h_c0)
    # h_co is hashable, so if it can be a key in a dict
    # then this test passes.
    dict_h_c0 = {h_c0: True}
    assert.contains(dict_h_c0, h_c0)

    ### Str and next ###
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

    c7 = count(5.0, 0.5)
    assert.eq(str(c7), "count(5.0, 0.5)")
    assert.eq(next(c7), 5.0)
    assert.eq(str(c7), "count(5.5, 0.5)")
    assert.eq(next(c7), 5.5)
    assert.eq(str(c7), "count(6.0, 0.5)")
    assert.eq(next(c7), 6.0)

    c8 = count(5, 1.0)
    assert.eq(str(c8), "count(5, 1.0)")

    ### NaNs ###
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

    ### Fails ###
    # Non-numeric arg fails.
    assert.fails(
        lambda: count("5", 1),
        'a number is required',
    )
    assert.fails(
        lambda: count(5, "1"),
        'a number is required',
    )
    assert.fails(
        lambda: count(None, 1),
        'a number is required',
    )
    assert.fails(
        lambda: count(5, None),
        'a number is required',
    )
    # Too many arg fails — should be handled by UnpackArgs but
    # check to be exhaustive.
    assert.fails(
        lambda: count(1, 2, 3),
        'count: got 3 arguments, want at most 2'
    )


def test_islice():

    assert.eq(
        list(islice([], 1)),
        []
    )
    assert.eq(
        list(islice(count(0, 5), 5)),
        [0, 5, 10, 15, 20]
    )
    assert.eq(
        list(islice(count(0, 5), 1, 5)),
        [5, 10, 15, 20]
    )
    assert.eq(
        list(islice(count(0, 5), 1, 5, 3)),
        [5, 20]
    )
    assert.eq(
        list(islice({'a': 0, 'b': 0, 'c': 0}, 3)),
        ['a', 'b', 'c']
    )
    assert.eq(
        list(islice([1, 2, 3], None)),
        [1, 2, 3]
    )

    # Check hashibility
    s0 = islice([1, 2, 3], 1)
    s1 = islice([1, 2, 3], 1)
    assert.true(s0 != s1)
    assert.true(s0 == s0)
    assert.true(s1 == s1)

    # Test according to Python docs; specifically testing
    # slice attr initialization and islice iteration.
    # See https://docs.python.org/3/library/itertools.html#itertools.islice
    s2 = islice("ABCDEFG".elems(), 2)
    assert.eq(
        "".join(list(s2)),
        "AB"
    )
    s3 = islice("ABCDEFG".elems(), 2, 4)
    assert.eq(
        "".join(list(s3)),
        "CD"
    )
    s4 = islice("ABCDEFG".elems(), 2, None)
    assert.eq(
        "".join(list(s4)),
        "CDEFG"
    )
    s5 = islice("ABCDEFG".elems(), 0, None, 2)
    assert.eq(
        "".join(list(s5)),
        "ACEG"
    )

    assert.fails(
        lambda: islice("asd".elems(), "a"),
        # fails uses match under the hood, which will use
        # regexp.MatchString, so need to use raw pattern
        # that MatchString would accept.
        r'expected int or None, got string',
    )
    assert.fails(
        lambda: islice("asd".elems(), 1, "a"),
        # fails uses match under the hood, which will use
        # regexp.MatchString, so need to use raw pattern
        # that MatchString would accept.
        r'expected int or None',
    )
    assert.fails(
        lambda: islice("asd".elems(), 1, None, "a"),
        # fails uses match under the hood, which will use
        # regexp.MatchString, so need to use raw pattern
        # that MatchString would accept.
        r'expected int or None',
    )
    assert.fails(
        lambda: islice("asd".elems(), -1),
        # fails uses match under the hood, which will use
        # regexp.MatchString, so need to use raw pattern
        # that MatchString would accept.
        r'expected non-negative values',
    )
    assert.fails(
        lambda: islice("asd".elems(), 1, None, -1),
        # fails uses match under the hood, which will use
        # regexp.MatchString, so need to use raw pattern
        # that MatchString would accept.
        r'expected non-negative values',
    )


test_count()
test_islice()
