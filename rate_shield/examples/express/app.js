import express from 'express'

const app = express()

const rateLimitCheck = async (req, res, next) => {
    const apiPath = req.baseUrl + req.path

    const headers = {
        'endpoint': apiPath,
        'ip': req.ip.replace('::ffff:', '')
    }

    try {
        const response = await fetch('http://127.0.0.1:8080/check-limit', {
            method: 'GET',
            headers: headers
        })

        if (response.status === 429) {
            res.status(429).json({
                error : 'TOO MANY REQUESTS'
            })
            return
        }

        if (response.status === 500) {
            res.status(500).json({
                error: 'INTERNAL SERVER ERROR'
            })
            return
        }

    } catch (e) {
        console.error('Error in rate limit check:', err)
        return res.status(500).json({
            error: 'Rate limit service unavailable'
        })
    }
    next()
}

app.use(rateLimitCheck)

app.get('/api/v1/process', (req, res) => {
    res.status(200).json({
        'success': true
    })
})


app.listen(3001, () => {
    console.log('server running on port 3001')
})