package main

func getCommSuffixI(s1 string) (commongSuffixIndex int) {
	commongSuffixIndex = dmp.DiffCommonSuffix(lasttext, s1)
	return commongSuffixIndex
}

func getCommPrefix(s1 string) int {
	commonPrefixI := dmp.DiffCommonPrefix(lasttext, s1)
	return commonPrefixI
}
