export function setCookie(
  name: string,
  value: string,
  daysToExpire: number
): void {
  const expirationDate = new Date();
  expirationDate.setTime(
    expirationDate.getTime() + daysToExpire * 24 * 60 * 60 * 1000
  );
  const expires = "expires=" + expirationDate.toUTCString();
  document.cookie = name + "=" + value + "; " + expires + "; path=/";
}

export function getCookie(name: string): string | null {
  const cookieName = name + "=";
  const cookies = document.cookie.split(";");
  for (let i = 0; i < cookies.length; i++) {
    let cookie = cookies[i].trim();
    if (cookie.indexOf(cookieName) === 0) {
      return cookie.substring(cookieName.length, cookie.length);
    }
  }
  return null;
}
