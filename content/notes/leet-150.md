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
    Initial Thoughts: Idea is to iterate through the array and maintain a separate index for the next position where a non-val element should be placed.
    Approach:
    1. Intialize Two Pointers (pointers hold the memory address of a variable)
        - i: this pointer will iterate over each element in the array
        - k or index: this pointer will keep track of the position where the next non-'val' element should be placed. Intially this pointer will be set to 0.
    2. Iterate Over the Array
        - For each element in 'nums', if the element is not equal to 'val', we place it at the 'k'th position and increment k by 'l'. This way all the non='val' elements are moved to the start of the array. 
    3. Result
        - After the iteration 'k' will be the count of elements in 'nums' that are not equal to 'val' and the first 'k' elements of 'nums' will be the elements that are not 'val'.
    
    Code: 
    ```
    def removeElement(nums, val):
    k = 0
    for i in range(len(nums)):
        if nums[i] != val:
            nums[k] = nums[i]
            k += 1
    return k

    ```

    **In computer science, an in-place algorithm is an algorithm that operates directly on the input data structure without requiring extra space proportional to the input size. In other words, it modifies the input in place, without creating a separate copy of the data structure. An algorithm which is not in-place is sometimes called not-in-place or out-of-place. -Wikipedia

    Time Complexity

    O(n): Where n is the length of nums. We are iterating through the array once.

    Space Complexity

    O(1): No extra space is used. The modifications are done in place.

3. Remove Duplicates from Sorted Array
    Approach:
    1. Handle Edge Cases: If the array is empty, return 0 as there are no unique elements
    2. Two Pointers
        i: this pointer iterates through the array starting from the second element
        k: this pointer marks the position where the next unique should be placed. It starts from the first element, assuming the first element is unique.
    3. Iterate and Check for Duplicates:
        Compare each element with the previous one. If they are different, it means we have found a unique element.
        Place this unique element at the 'k'th position and increment k.
    4. Return and Count:
        After the iteration, 'k' will be the count of unique elements, and the first 'k' elements of 'nums' will be the unique elements. 

    Code:
    ```
    def removeDuplicates(nums):
    if not nums:
        return 0

    k = 1  # Start from the second element
    for i in range(1, len(nums)):
        if nums[i] != nums[i - 1]:
            nums[k] = nums[i]
            k += 1
    return k

    ```

    Time Complexity

    O(n): Where n is the length of nums. The solution involves a single pass through the array.

    Space Complexity

    O(1): No additional space is used, as the operation is performed in-place.

    Loop Initialization: for i in range(1, len(nums))
        The loop starts with i = 1, not 0. This is because we compare each element with its previous element. Starting at 1 ensures there's always a "previous" element.
        The loop continues until i reaches the length of nums, iterating through the entire array.

    Condition Check: if nums[i] != nums[i - 1]
        In each iteration, we check if the current element (nums[i]) is different from the previous element (nums[i - 1]).
        Since the array is sorted, duplicate elements are adjacent. So, this check effectively identifies when a new unique element is encountered.

    Placing Unique Elements: nums[k] = nums[i]
        When a new unique element is found, we place it at the kth position in nums.
        This step gradually moves unique elements to the beginning of the array, overwriting duplicates.
        Note that if k and i are the same (which happens when there are no duplicates so far), this operation doesn't change the array but is still important for consistency.

    Incrementing k: k += 1
        After placing a unique element, we increment k.
        This increment ensures that the next unique element found will be placed in the next position in the array, maintaining the order of unique elements.

3. Remove Duplicates from Sorted Array II
    Approach:
    1. Handle Edge Cases: If the array has less than 3 elements, all elements can stay as they are, and we return the length of the array.
    2. Two Pointers:
        - i: this pointer iterates through the array starting from the second element
        - k: this pointer marks the position where the next element (either a unique element or a second occurrence) should be placed. It starts from the second element, assuming the first two elements can always stay. 
    3. Iterate and Check for Duplicates:
        Compare each element with the element two positions before. If they are different, it means the current element is either unique or a permissible second occurrence. 
    4. Return the Count:
        - After the iteration, k will be the count of elements after allowing each unique element ot appear at most twice. 

    Code: 
    ```
    def removeDuplicates(nums):
    if len(nums) < 3:
        return len(nums)

    k = 2  # Start from the third element
    for i in range(2, len(nums)):
        if nums[i] != nums[k - 2]:
            nums[k] = nums[i]
            k += 1
    return k
    
    ```

    Time Complexity

    O(n): Where n is the length of nums. The solution involves a single pass through the array.

    Space Complexity

    O(1): No additional space is used, as the operation is performed in-place.