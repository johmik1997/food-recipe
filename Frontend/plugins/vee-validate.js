import { defineRule, configure } from 'vee-validate'
import { required, email, min, confirmed } from '@vee-validate/rules'
import { defineNuxtPlugin } from '#app'

export default defineNuxtPlugin(() => {
  // Define validation rules
  defineRule('required', required)
  defineRule('email', email)
  defineRule('min', min)
  defineRule('confirmed', confirmed)

  // Configure vee-validate
  configure({
    validateOnInput: true,
    generateMessage: (context) => {
      const messages = {
        required: 'This field is required',
        email: 'Please enter a valid email',
        min: `This field must be at least ${context.rule.params} characters`,
        confirmed: 'Passwords do not match'
      }
      return messages[context.rule.name] || 'Invalid field'
    }
  })
})