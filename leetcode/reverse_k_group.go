package leetcode

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */

type ListNode struct {
	Val  int
	Next *ListNode
}

func reverseKGroup(head *ListNode, k int) *ListNode {
	if head == nil || k == 1 {
		return head
	}

	// Check if there are at least k nodes left
	node := head
	for i := 0; i < k; i++ {
		if node == nil {
			return head
		}
		node = node.Next
	}

	// Reverse k nodes
	var prev *ListNode
	curr := head
	for i := 0; i < k; i++ {
		next := curr.Next
		curr.Next = prev
		prev = curr
		curr = next
	}

	// head is now the tail of the reversed k nodes
	// curr is the next node to process
	head.Next = reverseKGroup(curr, k)

	return prev
}


server {
    listen 80;
    server_name cuberouter.ai www.cuberouter.ai;
    # Perform a 302 (temporary) redirect
    return 302 https://www.cuberouter.com$request_uri;
}