const express = require('express');
const path = require('path');
const hostelRoutes = require('./routes/hostelRoutes');
dotenv.config()
const app = express();
const PORT = process.env.PORT || 3000;

// Set view engine
app.set('view engine', 'ejs'); // or any other templating engine you prefer
app.set('views', path.join(__dirname, 'views'));

// Serve static files
app.use(express.static(path.join(__dirname, 'assets')));
app.use(express.static(path.join(__dirname, 'styles')));

// Use routes
app.use('/', routers/hostelRoutes);

app.listen(PORT, () => {
    console.log(`Server is running on http://localhost:${PORT}`);
});