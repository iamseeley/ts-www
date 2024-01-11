---
title: Leet 150
summary: Leetcode Problems for interview prep
tags: [interview prep]
date: Jan, 11 2024
draft: true
---

Language Choice: Python

1. Merge Sorted Array

    Step by step approach:
    1. Start from the end: Since 'nums1' has enough space to accommodate the elements of 'nums2' we start filling 'nums1' from the end. This way we don't have to shift elements around.
    2. Use three pointers: 
    - i: poinsts to the last element of the non-zero part of nums1 (m-1)
    - j: points to the last element of nums2 (n-1)
    - k:points to the last position in 'nums1' (m+n-1)
        Why? 
            i: This pointer starts at the end of the meaningful elements in nums1 (i.e., at index m-1). It represents the current element in nums1 that we are considering for merging.
            j: This pointer starts at the end of nums2 (i.e., at index n-1). It represents the current element in nums2 that we are considering for merging.
            k: This is the most critical pointer. It starts at the very end of nums1 (i.e., at index m+n-1, the last position available in nums1). This is where we place the next largest element during the merge process.
        Pointer: a pointer is a variable that stores the memory address of another variable.
        Pointers:
            Memory Address: A pointer holds the address of a variable, which allows for direct access and manipulation of the memory where the data is stored.
            Indirection: Pointers allow for the level of indirection. Through a pointer, you can access and change the value of the variable it points to.
            Dynamic Memory Allocation: In some languages, pointers are used to dynamically allocate memory on the heap.
    3. Compare and fill: We compare the elements at i and j and place the larger one at position k, then decrement the respective pointers.

    Code:
    ```
    def merge(nums1, m, nums2, n):
    i, j, k = m-1, n-1, m+n-1

    while i >= 0 and j >= 0:
        if nums1[i] > nums2[j]:
            nums1[k] = nums1[i]
            i -= 1
        else:
            nums1[k] = nums2[j]
            j -= 1
        k -= 1

    # If there are remaining elements in nums2
    while j >= 0:
        nums1[k] = nums2[j]
        j -= 1
        k -= 1
    ```

    Time Complexity

    O(m+n): Each element from both arrays is looked at once.

    Space Complexity

    O(1): No extra space is used; the merge is done in place.

2. Remove Element