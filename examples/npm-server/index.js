const express = require('express');
const app = express();

const PORT = process.env.PORT;

if (!PORT){
  console.log("The PORT variable is not set.");
  process.exit(1);
}

app.get("/", (req, res) => {
  res.status(200).json({ status: 'OK' });
});

app.get('/health', (req, res) => {
  res.status(200).json({ status: 'OK' });
});

// Start the server
app.listen(PORT, () => {
  console.log(`Server is running on port ${PORT}`);
});
