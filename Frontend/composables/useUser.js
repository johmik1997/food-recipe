export const useUser = () => {
return useState('user', () => {
  const userCookie = useCookie('user_data')
  const localUser = localStorage.getItem('user')
  return userCookie.value || (localUser ? JSON.parse(localUser) : null)
})
}