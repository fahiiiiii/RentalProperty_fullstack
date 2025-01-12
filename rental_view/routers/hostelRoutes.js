const express = require('express');
const router = express.Router();
const hostelController = require('../controllers/hostelController');

// Route to get the hostel listing
router.get('/hostel', hostelController.getHostelListing);

module.exports = router;