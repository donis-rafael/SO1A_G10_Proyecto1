const mongoose = require('mongoose')

const subscriberSchema = new mongoose.Schema({
  name: {
    type: String,
    required: true
  },
  location: {
    type: String,
    required: true
  },
  age: {
    type: Number,
    required: true 
  },
  infectedtype: {
    type: String,
    required: true 
  },
  state: {
    type: String,
    required: true 
  },
  origen: {
    type: String,
    required: true 
  }
})

module.exports = mongoose.model('Subscriber', subscriberSchema)