import React from 'react';

function Newsletter() {
  return (
    <section>
      <div class="columns is-paddingless">
        <div class="column is-8 is-offset-2 has-text-centered">
          <h3 class="title is-4 is-marginless">
            That sounds awesome, please keep me up to date!
          </h3>
          <form
            action="https://gmail.us20.list-manage.com/subscribe/post?u=cb2a947d7ab1a96ab7b6c8a54&amp;id=b2cbd3d806"
            method="post"
            id="mc-embedded-subscribe-form"
            name="mc-embedded-subscribe-form"
            class="validate"
            target="_blank"
            novalidate
          >
            <div class="field is-grouped">
              <p class="control is-expanded">
                <input
                  class="input"
                  type="email"
                  placeholder="Your email"
                  name="EMAIL"
                  id="mce-EMAIL"
                />
              </p>
              <p class="control">
                <button type="submit" class="button is-primary">
                  Subscribe
                </button>
              </p>
            </div>
            <div id="mce-responses" class="clear">
              <article class="message is-success">
                <div
                  id="mce-success-response"
                  class="message-body"
                >
                  Success!
                </div>
              </article>
            </div>
            <div aria-hidden="true">
              <input
                type="text"
                name="b_cb2a947d7ab1a96ab7b6c8a54_b2cbd3d806"
                tabindex="-1"
                value=""
              />
            </div>
          </form>
        </div>
      </div>
    </section>
  );
}

export default Newsletter;
