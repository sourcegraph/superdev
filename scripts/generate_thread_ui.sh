#!/bin/bash

# Script to run superdev with thread UI implementation prompts

# 1. Simple List View
echo "Running superdev command for Simple List View..."
superdev run Dockerfile --prompt "Implement a clean, minimal thread list view component that displays all threads in a tabular format. The table should have the following columns:
- Thread ID (truncated with copy button)
- Status (with color coding: yellow for processing, green for completed, red for error)
- Creation Time (formatted as relative time, e.g., '2 minutes ago')

The table should have a fixed header row, pagination for more than 10 entries, and responsive design that works on both desktop and mobile. Each row should be clickable to navigate to the thread detail page. Use the /threads API endpoint to fetch the list data and implement auto-refresh every 10 seconds.

Technical requirements:
- React functional component with hooks
- Fetch data from '/threads' endpoint
- Handle loading, error and empty states appropriately
- Sort by creation time (newest first) by default"

echo "Completed Simple List View"
echo "----------------------------------------"

# 2. Basic Status Filter Buttons
echo "Running superdev command for Basic Status Filter Buttons..."
superdev run Dockerfile --prompt "Create a simple filter control consisting of four buttons positioned above the thread list:
- All (default selected)
- Processing
- Completed
- Error

Each button should show the count of threads in that status category in a small badge. When a filter is active, apply a visual indicator to the selected button and filter the thread list to show only matching threads. The filtering should happen client-side after data is fetched to avoid additional API calls.

The filter state should be preserved in the URL query parameter to support bookmarking and sharing filtered views. Implement this feature using React state management and URL manipulation.

Technical requirements:
- React functional component with hooks
- Integrate with the thread list component
- Use URL parameters for filter state persistence
- Update counts dynamically as thread statuses change"

echo "Completed Basic Status Filter Buttons"
echo "----------------------------------------"

# 3. Thread Detail Page
echo "Running superdev command for Thread Detail Page..."
superdev run Dockerfile --prompt "Develop a thread detail page that shows comprehensive information about a specific thread. The page should include:
- Thread ID (with copy button)
- Status indicator 
- Creation timestamp (full date and time)
- Repository link (clickable)
- Full output text in a monospace font with proper whitespace preservation
- Error message (if status is 'error')
- Back button to return to the thread list

Use the /output?thread_id={id} endpoint to fetch the thread details. Implement proper error handling for non-existent thread IDs or network issues. The output text should be in a scrollable container with line numbers.

Technical requirements:
- React functional component with React Router integration
- Fetch data from '/output' endpoint with proper thread_id
- Implement loading state with skeleton UI
- Handle error states with appropriate user feedback
- Style the code output for readability"

echo "Completed Thread Detail Page"
echo "----------------------------------------"

# 4. Minimal Header with Thread Counts
echo "Running superdev command for Minimal Header with Thread Counts..."
superdev run Dockerfile --prompt "Create a simple application header that displays a summary of thread counts by status. The header should:
- Be fixed at the top of the application
- Show the total number of threads
- Display counts of threads by status (Processing/Completed/Error)
- Include a manual refresh button with loading indicator
- Show the last refresh timestamp

The header should be compact and use minimal space while providing useful information at a glance. The counts should update automatically when the thread list refreshes.

Technical requirements:
- React functional component that consumes thread data from a shared context
- Subscribe to thread data updates
- Responsive design that adapts to different screen sizes
- Accessible design with proper contrast and focus states"

echo "Completed Minimal Header with Thread Counts"
echo "----------------------------------------"

# 5. Single-Click Copy Functionality
echo "Running superdev command for Single-Click Copy Functionality..."
superdev run Dockerfile --prompt "Implement a copy-to-clipboard feature for thread IDs and output text. Add small copy icons next to:
- Thread IDs in the list view
- Thread ID in the detail view
- Repository link in the detail view
- A 'Copy All' button for the output text

When clicked, the icon should change briefly to indicate successful copying with a small toast notification. Implement this feature using the clipboard API with proper fallbacks for browsers that don't support it.

Technical requirements:
- Create a reusable CopyButton component
- Use the Clipboard API with proper fallbacks
- Provide visual feedback on copy success/failure
- Implement proper keyboard accessibility
- Ensure the component works across different browsers"

echo "Completed Single-Click Copy Functionality"
echo "----------------------------------------"

echo "All superdev commands completed"