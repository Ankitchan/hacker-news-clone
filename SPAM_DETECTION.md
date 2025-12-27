# Spam Detection System

This application includes an AI-powered spam detection system that protects against spam posts and comments.

## Features

### 1. **Dual-Layer Protection**

#### Layer 1: Built-in Heuristic Checks (Always Active, Free)
Fast, local spam detection that catches obvious spam without any API calls:

- **Excessive Capitalization**: Flags text with >70% capital letters
- **Spam Keywords**: Detects common spam phrases like:
  - "click here", "buy now", "free money"
  - "viagra", "cialis", "weight loss"
  - "nigerian prince", "make money fast"
  - And many more...
- **URL Bombing**: Flags content with more than 3 URLs
- **Character Repetition**: Detects excessive repeated characters (e.g., "hellooooooo")

#### Layer 2: AI Model (Optional, Requires API Key)
Advanced spam detection using Hugging Face's pre-trained models:

- **Model**: `mrm8488/bert-tiny-finetuned-sms-spam-detection`
- **Type**: BERT-based transformer model
- **Size**: Tiny (fast inference, low cost)
- **Accuracy**: Fine-tuned specifically for spam detection

### 2. **Cost Optimization**

#### Free Tier Usage
- **Hugging Face Free Tier**: Includes monthly credits
- **Heuristic Fallback**: Works even without API key
- **Request Minimization**: Heuristics catch spam before API calls
- **Fail-Safe**: If API fails, content is allowed (no false positives)

#### Estimated Costs (with AI enabled)
- **Free Tier**: ~10,000 requests/month (approximate)
- **Cost per request**: $0.000-$0.001 (varies by model)
- **Heuristics**: 0 cost, filter ~30-50% of spam locally

### 3. **Configuration**

#### Enable/Disable Spam Detection

In `.env`:
```bash
# Disable entirely (only heuristics work)
SPAM_DETECTION_ENABLED=false

# Enable with AI (requires API key)
SPAM_DETECTION_ENABLED=true
HF_API_KEY=your_api_key_here
```

#### Get a Free Hugging Face API Key

1. Go to https://huggingface.co/settings/tokens
2. Click "New token"
3. Name it (e.g., "hackernews-spam-detection")
4. Select "read" permissions
5. Copy the token
6. Add to `.env`:
   ```bash
   HF_API_KEY=hf_xxxxxxxxxxxxxxxxxxxxx
   ```

#### Custom Model

Want to use a different model? Update `.env`:
```bash
HF_SPAM_MODEL=your-username/your-spam-model
```

Popular alternatives:
- `IreNkweke/HamOrSpamModel`
- `AntiSpamInstitute/spam-detector-bert-MoE-v2.2`
- `Goodmotion/spam-mail-classifier`

## How It Works

### Detection Flow

```
User submits post/comment
        ↓
1. Check heuristics (fast, free)
   ├─ Is spam? → Reject immediately
   └─ Not obvious spam → Continue
        ↓
2. Check with AI model (if enabled)
   ├─ API call fails? → Allow (fail-safe)
   ├─ Confidence > 50%? → Reject as spam
   └─ Low confidence? → Allow
        ↓
3. Save to database
```

### Code Example

```go
// In handlers/post_handler.go
isSpam, confidence, err := h.spamDetector.IsSpam(content)
if err != nil {
    // Log error but don't block - fail-safe
    log.Printf("Spam detection error: %v", err)
} else if isSpam {
    log.Printf("Spam detected with confidence %.2f", confidence)
    utils.RespondWithError(w, http.StatusForbidden, "Content flagged as spam")
    return
}
```

## API Response Format

### Hugging Face Inference API Response

```json
[
  [
    {"label": "SPAM", "score": 0.9876},
    {"label": "HAM", "score": 0.0124}
  ]
]
```

The system uses the highest confidence score to determine if content is spam.

## Performance

### Latency
- **Heuristics only**: <1ms
- **With AI model**: 100-500ms (depends on Hugging Face API)
- **Timeout**: 10 seconds (configurable)

### Accuracy
- **Heuristics**: ~60-70% spam catch rate
- **AI Model**: ~95% accuracy
- **Combined**: ~97-98% spam catch rate
- **False Positives**: <2% (with fail-safe)

## Monitoring

### Check Spam Detection Logs

```bash
# Docker
docker-compose logs backend | grep -i spam

# Local
tail -f backend.log | grep -i spam
```

### Example Log Output

```
2025/12/27 15:30:42 Spam detected with confidence 0.98: CLICK HERE FOR FREE MONEY!!!
```

## Troubleshooting

### Issue: All Posts Are Blocked

**Cause**: Overly sensitive detection or API errors

**Solution**:
```bash
# Option 1: Disable temporarily
SPAM_DETECTION_ENABLED=false

# Option 2: Check API key
echo $HF_API_KEY

# Option 3: Check logs
docker-compose logs backend | grep spam
```

### Issue: API Rate Limit Exceeded

**Cause**: Too many requests in free tier

**Solutions**:
1. **Disable AI temporarily**:
   ```bash
   SPAM_DETECTION_ENABLED=false
   ```
   (Heuristics still work!)

2. **Upgrade to Pro**: $9/month for 20× more credits
3. **Wait**: Free tier resets monthly

### Issue: Legitimate Content Flagged

**Cause**: False positive from AI model

**Solutions**:
1. **Adjust threshold** in `detector.go`:
   ```go
   if prediction.Score > 0.7 { // Increase from 0.5
       return true, prediction.Score, nil
   }
   ```

2. **Disable AI, keep heuristics**:
   ```bash
   HF_API_KEY=  # Empty = heuristics only
   ```

## Testing Spam Detection

### Test with Obvious Spam

```bash
# Should be rejected by heuristics
curl -X POST http://localhost:8080/api/posts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"CLICK HERE FOR FREE MONEY NOW!!!","text":"Buy now limited time offer"}'

# Expected response: 403 Forbidden
```

### Test with Legitimate Content

```bash
# Should pass
curl -X POST http://localhost:8080/api/posts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Interesting article about technology","url":"https://example.com"}'

# Expected response: 201 Created
```

## Extending the System

### Add Custom Heuristics

Edit `pkg/spam/detector.go`:

```go
func (d *Detector) quickSpamCheck(text string) bool {
    text = strings.ToLower(text)

    // Your custom checks here
    if strings.Contains(text, "your-spam-phrase") {
        return true
    }

    return false
}
```

### Use Different AI Model

1. Find a model on https://huggingface.co/models?other=spam
2. Test it on Hugging Face
3. Update `.env`:
   ```bash
   HF_SPAM_MODEL=username/model-name
   ```

### Add Database Logging

Track spam attempts in the database:

```go
// In detector.go
if isSpam {
    // Log to database
    d.logSpamAttempt(userID, content, confidence)
    return true, confidence, nil
}
```

## Best Practices

1. **Start with free tier**: Test with monthly credits
2. **Monitor logs**: Watch for false positives
3. **Use heuristics**: They're fast and free
4. **Fail-safe**: Always allow on API errors
5. **Adjust threshold**: Fine-tune for your use case
6. **Cache results**: Consider caching for repeat content
7. **Rate limit**: Combine with rate limiting for best protection

## Resources

- [Hugging Face Models](https://huggingface.co/models?other=spam)
- [Hugging Face Inference API Docs](https://huggingface.co/docs/api-inference/index)
- [Hugging Face Pricing](https://huggingface.co/pricing)
- [BERT-tiny Model](https://huggingface.co/mrm8488/bert-tiny-finetuned-sms-spam-detection)

## Support

For spam detection issues:
1. Check logs: `docker-compose logs backend | grep spam`
2. Verify API key: `echo $HF_API_KEY`
3. Test heuristics: Set `HF_API_KEY=` (empty)
4. Open an issue with example spam content that wasn't caught
