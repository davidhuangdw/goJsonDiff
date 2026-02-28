package diff

func longestCommonSubHashes(fromHashes, toHashes []*HashTree) []uint64 {
	fr, to := make([]uint64, len(fromHashes)), make([]uint64, len(toHashes))
	for i, h := range fromHashes {
		fr[i] = h.Hash
	}
	for i, h := range toHashes {
		to[i] = h.Hash
	}
	return longestCommonSublist(fr, to)
}

func longestCommonSublist[T comparable](a, b []T) []T {
	n, m := len(a), len(b)
	longest := make([][]int, n+1)
	for i := range longest {
		longest[i] = make([]int, m+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if a[i] == b[j] {
				longest[i][j] = 1 + longest[i+1][j+1]
			} else {
				longest[i][j] = max(longest[i][j+1], longest[i+1][j])
			}
		}
	}

	common := make([]T, 0)
	i, j := 0, 0
	for i < n && j < m {
		if a[i] == b[j] {
			common = append(common, a[i])
			i++
			j++
		} else if longest[i][j] == longest[i+1][j] {
			i++ // prefer to first remove than add
		} else {
			j++
		}
	}
	return common
}
