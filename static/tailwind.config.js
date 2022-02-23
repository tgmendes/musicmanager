module.exports = {
    content: ["./templates/**/*.{html,js}"],
    theme: {
        fontFamily: {
            body: ['Raleway'],
            display: ['Prompt']
        },
        extend: {
            colors: {
                primary: {
                    100: '#F2ABB8',
                    200: '#EC8294',
                    300: '#E65870',
                    400: '#DF2E4D',
                    500: '#D90429'
                },
                secondary: {
                    100: '#B8B9C0',
                    200: '#9596A1',
                    300: '#727381',
                    400: '#4E5062',
                    500: '#2B2D42',
                },
            },
        },
    },
    plugins: [],
}