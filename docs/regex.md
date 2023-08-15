# Regular Expressions

Many of the checks in `aeacus` use regular expression (regex) strings as input. This may seem inconvenient if you want
to score something simple, but we think it significantly increases the overall quality of checks. Each regex is applied
to each line of the input file, so currently, no multi-line regexes are currently possible.

The checks that are specifically supported are `CommandContainsRegex`, `DirContainsRegex`, and `FileContainsRegex`.
Please note that you **must** append `Regex` to the end for the check to use regular expressions.

> We're using the Golang Regular Expression package ([documentation here](https://godocs.io/regexp)). It uses RE2
> syntax, which is also generally the same as Perl, Python, and other languages.

If you're unfamiliar, a 'regex' is just a way of describing a pattern of text. Let's say I was trying to score this
in `/etc/apt/apt.conf.d/*`:

```
APT::Periodic::Update-Package-Lists "1";
```

The most simple regexes work the same way that normal CTRL-f searches work. It just matches what it is. A valid regex to
score this would be `APT::Periodic::Update-Package-Lists "1";`, since none of those characters mean anything special to
regex. With normal substring searching (like CTRL-f), that's the most specific we would be able to be.

But, what about this?

```
APT::Periodic::Update-Package-Lists  "1";
```

Notice that there are now two spaces before the `"1"`. That's a bummer, because our config is still valid to Ubuntu's
software updater, but it's not scored as correct. We need at least one space between those two, so we'll
do `APT::Periodic::Update-Package-Lists\s+"1";`, where `\s` means any whitespace, and `+` means 'at least one.'

Similarly, what about all of these?

```
APT::Periodic::Update-Package-Lists  "1";
APT::Periodic::Update-Package-Lists "1" ;
  APT::Periodic::Update-Package-Lists "1";
```

Our new regex, to match all of those, would be `\s*APT::Periodic::Update-Package-Lists\s+"1"\s*;\s*`. `*` means 'any
amount of the preceding token, including none.' So this will match any amount of whitespace (including no whitespace).
Which would work much better.

But, what about this?

```
# APT::Periodic::Update-Package-Lists "1";
```

That ruins everything. It would score as correct even though it's commented out. But, it's an easy fix. With the
regex `^\s*APT::Periodic::Update-Package-Lists\s+"1"\s*;\s*$`, where we added `^` for 'start of the line' and `$` for '
end of the line', it will only match if there's nothing except whitespace before and after the directive.

But, what about this?

```
APT::PERIODIC::UPDATE-PACKAGE-LISTS "1";
```

Believe it or not, the apt configs appear to be case-insensitive. So we modify the expression to be case insensitive
with `(?i)`: `(?i)^\s*APT::Periodic::Update-Package-Lists\s+"1"\s*;\s*$`.

As far as I know, this is as correct as we can get it.

Thinking about the edge cases and correct grammar for scoring these directives is very important and makes a big
difference in scoring robustness, which is why we use regexes for many checks. It can take a lot of practice to get a
working expression, and making mistakes is very common. If you want to test your expression interactively, you can use
something like [debuggex](https://debuggex.com) or [regex101](https://regex101.com).
