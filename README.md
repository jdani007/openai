## OPENAI-CHAT

Based on Python implementation from Linkedin course [Building a Project with the ChatGPT API](https://www.linkedin.com/learning/building-a-project-with-the-chatgpt-api/challenge-generate-tweet-content-using-chat-completion?resume=false&u=2288114)

[Create chat completion](https://platform.openai.com/docs/api-reference/chat/create)

[Create image](https://platform.openai.com/docs/api-reference/images/create)

```
User: You are a friendly AI assistant that helps compose professional-sounding tweets for Twitter that often go viral based on a website I provide. You will provide a summary of the website in 50 words or less.

Assistant: Sure, I'd be happy to help! Please provide me with the website you'd like a summary of.
```

```
User: https://learning.linkedin.com

Assistant: LinkedIn Learning provides online courses to help professionals enhance their skills and advance their careers. Topics cover a wide range of subjects, from business and technology to creative skills and personal development. It's a valuable resource for continuous learning and professional growth in today's competitive job market.
```


```
User: Can you shorten the summary?

Assistant: LinkedIn Learning offers a variety of online courses for professionals to improve skills and boost careers. With topics ranging from business to personal development, it's a valuable tool for continuous learning and advancing in the competitive job market.
```

```
User: Can you describe an image that represents the summary?

Assistant: Imagine a vibrant collage of diverse professionals engaged in online courses on laptops and tablets against a backdrop of colorful educational icons symbolizing a variety of skills. This visual captures the essence of continuous learning, career advancement, and personal growth that LinkedIn Learning offers.
```

```
User: generate image

Assistant: Generating image............
```
Image will be generated and opened via web browser (I use a mac).
![generated image](generated-image.png)