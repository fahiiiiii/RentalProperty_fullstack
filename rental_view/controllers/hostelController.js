const hostelData = require('../models/hostelModel');

exports.getHostelListing = (req, res) => {
    const hostelInfo = hostelData.getHostelInfo(); // Assuming this fetches the hostel data
    res.render('hostel-listing', { hostel: hostelInfo });
};