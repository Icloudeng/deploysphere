export const wait = (sec = 1) =>
  new Promise((resolve) => setTimeout(resolve, sec * 1000));

export const log = (text: string) => console.log(`%%%${text}%%%`);
