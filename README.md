# MusicSharingAPI

 User creation:
   - Signup
   - Login
   - Logout
   - Password encryption and decryption
   - Unique usernames (only one person can have a username)
   - Gorilla session management for easy session handling

User interaction:
  - Users can follow and unfollow other users (following only allowed if a user is not already followed)
  - Following, followers, and posts count
  - Users can share posts with their followers
  - A feed-style homepage with all of a user's following users' posts

Post interaction:
  - Users can like/unlike posts (only one like per user per post)
  - Users can Delete posts as long as they are the owner of the post
  - Posts feature an author which is automatically set to the username of the current user, and a like count
  - Each user has a slice of posts to view previously liked posts
