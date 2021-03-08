//when generating paths from the top level, run each searchBranch in its
//own goroutine, so there should be ten in the base case, all doing reads
//on the database; this should be fine, since only one request can be served
//by each context at any time and creation and learning are separate flows --
//creation is strictly read-only


//scoring logic:
//All productions start with a score of 0
//each test adds or removes points (typically 1 or 2) depending on how well the production
//satisfies its requirements:
//each primary keyword is worth one point; each secondary worth one
//failing to meet a minimum target length will deduct a point
//  slightly exceeding the target will award one, but greatly exceeding it will award nothing
//a point will be deducted for every repetition of the same token above two counts

//any production with a positive score is a response candidate and will be formatted
//productions are grouped by score and returned is descending order



//when producing, do N forward walks from the keyword and N reverse walks,
//then, for each of the paths that come back (probably grouped by common
//pattern), do a reverse-walk that looks at the full n-gram pattern and
//combine those, rather than the two-start-from-keyword MegaHAL approach.

//if there are no viable chains after scoring, then do N forward and reverse
//walks from the start and end positions, score them, and return that

